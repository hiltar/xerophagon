# Xerophagon

**Xerophagon** is a lightweight, retro-themed web application for tracking fasting sessions. Built with Go, it runs efficiently on a Raspberry Pi Zero 2W using Alpine Linux** and **Podman/Docker** for containerization. The app provides a user-friendly interface to start, monitor, and end fasts, visualize fasting stages, and review fasting history with a vibrant, 1980s-inspired retro aesthetic.

![Profile](https://github.com/user-attachments/assets/4071b3a1-fb56-4e25-8a07-2ba948375711)  
![Fasting stages](https://github.com/user-attachments/assets/3daac0e9-2d89-4084-8273-7c9b6eba2424)  
![Fasting](https://github.com/user-attachments/assets/0e096b5e-cd57-4def-825f-9a3c8b8c4c46)  

## Features

- **Fasting Management**:
  - Start a new fast with an optional goal (0 hours for open-ended fasts).
  - End fasts and log duration in hours.
  - Real-time timer displaying elapsed fasting time.
- **Fasting Stages Visualization**:
  - Displays a dot grid representing fasting hours (14px dots, 12 per row), colored by physiological stages (e.g., Ketosis, Autophagy).
  - Current stage container highlights the active stage with vibrant colors and descriptions.
  - Modal view of all stages with detailed explanations.
- **Fasting History**:
  - Paginated fasting history.
  - Stores start time and duration for each completed fast.
- **Retro-Themed UI**:
  - 1980s-inspired design with neon cyan/pink colors, and pixelated 'VCR OSD Mono' font.
  - Responsive layout using **Bulma CSS**, optimized for mobile and desktop.
- **Lightweight and Offline**:
  - Runs on a Raspberry Pi Zero 2W with minimal resources.
  - Persistent data storage in a JSON file.

## Technical Details

### Architecture

- **Backend**:
  - Written in **Go 1.24** using the standard library's `net/http` package.
  - RESTful API-like endpoints (`/profile`, `/fasting`, `/start_fast`, `/end_fast`).
  - Data stored in a JSON file (`/data/data.json`) with a simple `AppData` struct:
    ```go
    type AppData struct {
        CurrentFast    *Fast  `json:"current_fast"`
        FastingHistory []Fast `json:"fasting_history"`
    }
    type Fast struct {
        StartTime     time.Time `json:"start_time"`
        GoalHours     int       `json:"goal_hours,omitempty"`
        DurationHours float64   `json:"duration_hours,omitempty"`
    }
    ```
  - Server-side pagination for fasting history (5 entries per page) using query parameters (e.g., `/profile?page=2`).
  - Template rendering with Go's `html/template` package.

- **Frontend**:
  - **HTML Templates**: `profile.html` and `fasting.html` for Profile and Fasting tabs.
  - **CSS**: Custom retro-themed styles in `style.css` with **Bulma CSS** for layout and components.
    - Retro aesthetic: neon colors, pixelated borders.
    - Local 'VCR OSD Mono' font via `@font-face` for offline use.
  - **JavaScript**: Minimal client-side logic in `script.js` for:
    - Real-time timer updates (every second).
    - Dot grid coloring based on fasting stages.
    - Current stage updates with vibrant backgrounds.
    - Modal interactions for stage details.

- **Containerization**:
  - Built with **Podman/Docker** using a multi-stage `Dockerfile`:
    - Builder stage: `golang:1.24-alpine` compiles the binary with `CGO_ENABLED=0` and compresses it using `upx`.
    - Final stage: `scratch` image for minimal footprint (~3 MB).
  - Persistent data volume at `/data` for `data.json`.
  - Runs as non-root user (1000:1000) for security.

### How It Works

1. **Starting a Fast**:
   - User navigates to the Fasting tab and submits a goal (hours, optional; 0 for open-ended).
   - `startFastHandler` validates input, allowing 0 or positive integers, and stores the fast in `data.json`.
   - Invalid inputs (e.g., negative or non-numeric) redirect with a custom error message ("Goal must be a non-negative number").

2. **Monitoring a Fast**:
   - Profile tab displays:
     - **Dot Grid**: 14px dots (12 per row) up to the goal hours, colored by stage (e.g., #1e90ff for Autophagy at 24+ hours). Current hour highlighted with a neon cyan border and pulse animation.
     - **Current Stage**: Box with stage-specific background, white text, and shadow for readability.
     - **Timer**: Updates every second via JavaScript.
   - Stages are defined in `script.js` (e.g., Ketosis at 12 hours, Autophagy at 24 hours).
   - A 0-hour goal hides the dot grid but shows the timer and stage.

3. **Ending a Fast**:
   - User submits the "End Fast" form in the Fasting tab.
   - `endFastHandler` calculates the duration and appends the fast to `FastingHistory` in `data.json`.

4. **Fasting History**:
   - Displayed in the Profile tab, paginated (5 entries per page).
   - Server-side logic in `profileHandler` slices the history based on the `page` query parameter.
   - Navigation via Previous/Next buttons styled with retro neon colors.

5. **Retro UI**:
   - **Background**: Dark gradient with scanline overlay (base64 PNG).
   - **Typography**: 'VCR OSD Mono' font, served locally from `static/fonts/`.
   - **Colors**: Neon cyan (#00ffcc) for text, pink (#ff00ff) for buttons, with vibrant stage colors (e.g., #32cd32 for Ketosis).
   - **Containers**: Retro boxes with neon borders and shadows.

## Requirements

- **Hardware**: Raspberry Pi Zero 2W (or similar low-power device).
- **OS**: Alpine Linux (or any Linux distro with Podman).
- **Software**:
  - **Docker/Podman** for containerization.
    

## Installation

### Prerequisites

1. **Build an executable**
   ```
   go mod init xerophagon
   go mod tidy
   CGO_ENABLED=0 GOOS=linux go build -a -ldflags="-s -w" -o xerophagon

   chmod a+x xerophagon
   ./xerophagon
   ```
   
2. **Build with Docker/Podman**
   ```
   sudo docker build -t xerophagon:v0.0.1 .
   sudo docker run -dt -p 5000:5000 --name xerophagon -v $(pwd)/data:/data xerophagon:v0.0.1
   ```

3. **Build for ARM architecture**  
   ![Instructions in arm branch](https://github.com/hiltar/xerophagon/blob/arm/README.md)
