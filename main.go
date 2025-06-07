package main

import (
    "embed"
    "encoding/json"
    "fmt"
    "html/template"
    "io/fs"
    "net/http"
    "os"
    "path/filepath"
    "strconv"
    "time"
)

//go:embed static templates
var content embed.FS

// Fast represents a fasting record
type Fast struct {
    StartTime     time.Time `json:"start_time"`
    GoalHours     int       `json:"goal_hours,omitempty"`
    DurationHours float64   `json:"duration_hours,omitempty"`
}

// AppData holds the application's state
type AppData struct {
    CurrentFast    *Fast  `json:"current_fast"`
    FastingHistory []Fast `json:"fasting_history"`
}

// TemplateData includes pagination data
type TemplateData struct {
    CurrentFast    *Fast
    FastingHistory []Fast
    CurrentPage    int
    TotalPages     int
    HasPrev        bool
    HasNext        bool
}

const dataFile = "./data/data.json"

// loadData reads the JSON data file
func loadData() (*AppData, error) {
    if _, err := os.Stat(dataFile); os.IsNotExist(err) {
        return &AppData{FastingHistory: []Fast{}}, nil
    }
    file, err := os.Open(dataFile)
    if err != nil {
        return nil, fmt.Errorf("opening data file: %w", err)
    }
    defer file.Close()
    var data AppData
    err = json.NewDecoder(file).Decode(&data)
    if err != nil {
        return nil, fmt.Errorf("decoding data file: %w", err)
    }
    return &data, nil
}

// saveData writes the application data to the JSON file
func saveData(data *AppData) error {
    dir := filepath.Dir(dataFile)
    if err := os.MkdirAll(dir, 0755); err != nil {
        return fmt.Errorf("creating data directory: %w", err)
    }
    file, err := os.Create(dataFile)
    if err != nil {
        return fmt.Errorf("creating data file: %w", err)
    }
    defer file.Close()
    encoder := json.NewEncoder(file)
    encoder.SetIndent("", "  ")
    if err := encoder.Encode(data); err != nil {
        return fmt.Errorf("encoding data: %w", err)
    }
    return nil
}

// homeHandler redirects to the profile page
func homeHandler(w http.ResponseWriter, r *http.Request) {
    http.Redirect(w, r, "/profile", http.StatusFound)
}

// profileHandler renders the Profile tab with pagination
func profileHandler(w http.ResponseWriter, r *http.Request) {
    data, err := loadData()
    if err != nil {
        http.Error(w, fmt.Sprintf("Error loading data: %v", err), http.StatusInternalServerError)
        return
    }

    const itemsPerPage = 5
    pageStr := r.URL.Query().Get("page")
    page, err := strconv.Atoi(pageStr)
    if err != nil || page < 1 {
        page = 1
    }
    totalItems := len(data.FastingHistory)
    totalPages := (totalItems + itemsPerPage - 1) / itemsPerPage
    if totalPages == 0 {
        totalPages = 1
    }
    start := (page - 1) * itemsPerPage
    end := start + itemsPerPage
    if start > totalItems {
        start = totalItems
    }
    if end > totalItems {
        end = totalItems
    }
    paginatedHistory := data.FastingHistory[start:end]

    templateData := TemplateData{
        CurrentFast:    data.CurrentFast,
        FastingHistory: paginatedHistory,
        CurrentPage:    page,
        TotalPages:     totalPages,
        HasPrev:        page > 1,
        HasNext:        page < totalPages,
    }

    tmpl, err := template.New("profile.html").Funcs(template.FuncMap{
        "seq": func(start, end int) []int {
            if end < start {
                return []int{}
            }
            s := make([]int, end-start+1)
            for i := range s {
                s[i] = start + i
            }
            return s
        },
        "mod": func(a, b int) int {
            if b == 0 {
                return 0
            }
            return a % b
        },
        "add": func(a, b int) int {
            return a + b
        },
    }).ParseFS(content, "templates/profile.html")
    if err != nil {
        http.Error(w, fmt.Sprintf("Error parsing template: %v", err), http.StatusInternalServerError)
        return
    }
    err = tmpl.Execute(w, templateData)
    if err != nil {
        http.Error(w, fmt.Sprintf("Error executing template: %v", err), http.StatusInternalServerError)
    }
}

// fastingHandler renders the Fasting tab
func fastingHandler(w http.ResponseWriter, r *http.Request) {
    data, err := loadData()
    if err != nil {
        http.Error(w, fmt.Sprintf("Error loading data: %v", err), http.StatusInternalServerError)
        return
    }
    t, err := template.ParseFS(content, "templates/fasting.html")
    if err != nil {
        http.Error(w, fmt.Sprintf("Error parsing template: %v", err), http.StatusInternalServerError)
        return
    }
    type FastingTemplateData struct {
        *AppData
        ErrorMessage string
    }
    templateData := FastingTemplateData{
        AppData:      data,
        ErrorMessage: r.URL.Query().Get("error"),
    }
    err = t.Execute(w, templateData)
    if err != nil {
        http.Error(w, fmt.Sprintf("Error executing template: %v", err), http.StatusInternalServerError)
    }
}

// startFastHandler handles starting a new fast
func startFastHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }
    err := r.ParseForm()
    if err != nil {
        http.Redirect(w, r, "/fasting?error=Error parsing form", http.StatusSeeOther)
        return
    }
    goalStr := r.FormValue("goal")
    var goal int
    if goalStr == "" {
        goal = 0
    } else {
        goal, err = strconv.Atoi(goalStr)
        if err != nil || goal < 0 {
            http.Redirect(w, r, "/fasting?error=Goal must be a non-negative number", http.StatusSeeOther)
            return
        }
    }
    data, err := loadData()
    if err != nil {
        http.Redirect(w, r, "/fasting?error=Error loading data", http.StatusSeeOther)
        return
    }
    if data.CurrentFast != nil {
        http.Redirect(w, r, "/fasting?error=A fast is already in progress", http.StatusSeeOther)
        return
    }
    data.CurrentFast = &Fast{
        StartTime: time.Now(),
        GoalHours: goal,
    }
    err = saveData(data)
    if err != nil {
        http.Redirect(w, r, fmt.Sprintf("/fasting?error=Error saving data: %v", err), http.StatusSeeOther)
        return
    }
    http.Redirect(w, r, "/fasting", http.StatusSeeOther)
}

// endFastHandler handles ending the current fast
func endFastHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }
    data, err := loadData()
    if err != nil {
        http.Error(w, fmt.Sprintf("Error loading data: %v", err), http.StatusInternalServerError)
        return
    }
    if data.CurrentFast == nil {
        http.Error(w, "No fast in progress", http.StatusBadRequest)
        return
    }
    duration := time.Since(data.CurrentFast.StartTime).Hours()
    completedFast := Fast{
        StartTime:     data.CurrentFast.StartTime,
        GoalHours:     data.CurrentFast.GoalHours,
        DurationHours: duration,
    }
    data.FastingHistory = append([]Fast{completedFast}, data.FastingHistory...)
    data.CurrentFast = nil
    err = saveData(data)
    if err != nil {
        http.Error(w, fmt.Sprintf("Error saving data: %v", err), http.StatusInternalServerError)
        return
    }
    http.Redirect(w, r, "/fasting", http.StatusSeeOther)
}

func main() {
    // Serve static files from embedded FS
    staticFS, err := fs.Sub(content, "static")
    if err != nil {
        fmt.Println("Error creating static FS:", err)
        os.Exit(1)
    }
    http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.FS(staticFS))))

    // Define routes
    http.HandleFunc("/", homeHandler)
    http.HandleFunc("/profile", profileHandler)
    http.HandleFunc("/fasting", fastingHandler)
    http.HandleFunc("/start_fast", startFastHandler)
    http.HandleFunc("/end_fast", endFastHandler)

    // Start the server
    fmt.Println("Xerophagon starting on :5000...")
    err = http.ListenAndServe(":5000", nil)
    if err != nil {
        fmt.Printf("Server failed: %v\n", err)
    }
}
