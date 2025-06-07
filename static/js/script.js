// Define fasting stages with colors
const stages = [
    { hours: 0, name: "Blood sugar is rising", description: "The body is digesting and absorbing nutrients. Blood sugar and insulin levels are elevated", class: "stage-0" },
    { hours: 4, name: "Blood sugar drops", description: "Blood sugar and insulin levels begin to drop, and the body starts to use stored glycogen for energy", class: "stage-4" },
    { hours: 8, name: "Fat burning phase", description: "Glycogen stores are depleted and the body begins to shift towards using fat for energy", class: "stage-8" },
    { hours: 12, name: "Ketosis phase", description: "The body enters ketosis, where fat stores are converted into ketones, providing an alternative energy source", class: "stage-12" },
    { hours: 18, name: "Higher ketosis phase", description: "Higher ketone levels signal the body to ramp up stress busting pathways to reduce inflammation and repair DNA damage", class: "stage-18" },
    { hours: 24, name: "Autophagy phase", description: "The body relies heavily on ketones for energy. Cellular repair and detoxification processes are at their peak", class: "stage-24" },
    { hours: 36, name: "Enhanced autophagy phase", description: "Body begins to recycle damaged cells and misfolded proteins", class: "stage-36" },
    { hours: 48, name: "Growth hormone peak phase", description: "Growth hormones reach the highest level (500%)", class: "stage-48" },
    { hours: 54, name: "Low insulin level phase", description: "Insulin levels drop to the lowest point. The body becomes more insulin sensitive", class: "stage-54" },
    { hours: 72, name: "Immune cell reformation phase", description: "Breakdown of old immune cells and generation of new ones", class: "stage-72" },
    { hours: 96, name: "Metabolic shift phase", description: "The body transitions to full ketosis, utilizing fatty acids and ketone bodies as primary energy sources", class: "stage-96" },
    { hours: 120, name: "Multi-Organ response phase", description: "Significant changes in protein levels across multiple organs, indicating systemic adaptation to prolonged fasting", class: "stage-120" },
    { hours: 144, name: "Protein level phase", description: "Systemic changes are occurring in protein levels, indicating widespread physiological adaptations", class: "stage-144" }
];

function updateUI() {
    const currentFastDiv = document.getElementById("current-fast");
    if (!currentFastDiv) return;

    const startTime = parseInt(currentFastDiv.dataset.startTime) * 1000; // Convert to milliseconds
    const now = Date.now();
    const elapsedMs = now - startTime;
    const elapsedSeconds = Math.floor(elapsedMs / 1000);
    const hours = Math.floor(elapsedSeconds / 3600);
    const minutes = Math.floor((elapsedSeconds % 3600) / 60);
    const seconds = elapsedSeconds % 60;
    const timerText = `${hours.toString().padStart(2, "0")}:${minutes
        .toString()
        .padStart(2, "0")}:${seconds.toString().padStart(2, "0")}`;
    document.getElementById("timer").textContent = timerText;

    // Update dot grid
    const dots = document.querySelectorAll(".dot");
    const elapsedHours = Math.floor(elapsedMs / 3600000);
    dots.forEach((dot) => {
        const hour = parseInt(dot.dataset.hour);
        // Find the stage for this hour
        let stageClass = "";
        for (let i = stages.length - 1; i >= 0; i--) {
            if (hour >= stages[i].hours) {
                stageClass = stages[i].class;
                break;
            }
        }
        // Apply stage color and filled status
        dot.className = `dot ${stageClass}`;
        if (hour <= elapsedHours) {
            dot.classList.add("filled");
        } else {
            dot.classList.remove("filled");
        }
        // Highlight current hour
        if (hour == elapsedHours) {
            dot.classList.add("current-hour");
        } else {
            dot.classList.remove("current-hour");
        }
    });

    // Update current stage
    const currentStageDiv = document.getElementById("current-stage");
    let currentStage = stages[0];
    for (let i = stages.length - 1; i >= 0; i--) {
        if (elapsedHours >= stages[i].hours) {
            currentStage = stages[i];
            break;
        }
    }
    currentStageDiv.innerHTML = `
        <strong>${currentStage.name}</strong>
        <p>${currentStage.description}</p>
    `;
    currentStageDiv.className = `box retro-box ${currentStage.class}`;

    // Populate stages modal
    const stagesList = document.getElementById("stages-list");
    stagesList.innerHTML = stages.map(stage => `
        <div class="stage-item ${stage.class}">
            <strong>${stage.name} (from ${stage.hours} hours)</strong>
            <p>${stage.description}</p>
        </div>
    `).join("");

    // Handle modal
    const modal = document.getElementById("stages-modal");
    const viewStagesBtn = document.getElementById("view-stages-btn");
    const closeButton = document.getElementById("modal-close-btn");
    const modalBackground = document.querySelector(".modal-background");
    viewStagesBtn.onclick = () => modal.classList.add("is-active");
    closeButton.onclick = () => modal.classList.remove("is-active");
    modalBackground.onclick = () => modal.classList.remove("is-active");
}

// Update every second
setInterval(updateUI, 1000);

// Initial call to avoid delay
updateUI();
