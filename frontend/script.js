const API_URL = "http://localhost:8080/workouts";

// Render the chart
function renderChart(data, exercise) {
  const ctx = document.getElementById(`chart-${exercise}`).getContext("2d");

  // Create a new chart
  new Chart(ctx, {
    type: "line",
    data: {
      labels: data.labels, // Dates
      datasets: [{
        label: `Workout Volume for ${exercise}`,
        data: data.volumes, // Volumes
        borderColor: "blue",
        backgroundColor: "rgba(0, 0, 255, 0.1)",
        fill: true,
      }]
    },
    options: {
      responsive: true,
      scales: {
        x: {
          title: {
            display: true,
            text: "Date"
          }
        },
        y: {
          title: {
            display: true,
            text: "Volume"
          }
        }
      }
    }
  });
}

// Fetch and display workouts
async function fetchWorkouts() {
  try {
    const response = await fetch(API_URL);
    if (!response.ok) {
      throw new Error(`Failed to fetch workouts: ${response.statusText}`);
    }
    const workouts = await response.json();
    const tableBody = document.querySelector("#workouts-table tbody");
    tableBody.innerHTML = ""; // Clear existing rows

    // Prepare data for the chart
    const groupedData = {};
    workouts.forEach(workout => {
      const volume = workout.sets.reduce((sum, set) => sum + set.weight * set.reps, 0);
      const row = `
        <tr>
          <td>${workout.date}</td>
          <td>${workout.exercise}</td>
          <td>${volume}</td>
        </tr>
      `;
      tableBody.innerHTML += row;

      if (!groupedData[workout.exercise]) {
        groupedData[workout.exercise] = { labels: [], volumes: [] };
      }
      groupedData[workout.exercise].labels.push(workout.date);
      groupedData[workout.exercise].volumes.push(volume);
    });

    // Find the full range of dates for alignment
    const allDates = Object.values(groupedData).flatMap(data => data.labels);
    const uniqueDates = Array.from(new Set(allDates)).sort();

    // Render a chart for each exercise
    const chartContainer = document.getElementById("visualization-section");
    chartContainer.innerHTML = ""; // Clear existing charts

    Object.entries(groupedData).forEach(([exercise, data]) => {
      // Align the data with the full range of dates
      const alignedVolumes = uniqueDates.map(date => {
        const index = data.labels.indexOf(date);
        return index !== -1 ? data.volumes[index] : 0;
      });

      // Create a canvas for the chart
      const canvas = document.createElement("canvas");
      canvas.id = `chart-${exercise}`;
      chartContainer.appendChild(canvas);

      // Render the chart
      renderChart({ labels: uniqueDates, volumes: alignedVolumes }, exercise);
    });
  } catch (error) {
    console.error("Error fetching workouts:", error);
  }
}

// Submit new workout
document.querySelector("#workout-form").addEventListener("submit", async (e) => {
  e.preventDefault();

  const date = document.querySelector("#date").value;
  const exercise = document.querySelector("#exercise").value;
  const sets = Array.from(document.querySelectorAll(".set")).map(set => ({
    weight: parseInt(set.querySelector(".weight").value),
    reps: parseInt(set.querySelector(".reps").value),
  }));

  const response = await fetch(API_URL, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ date, exercise, sets }),
  });

  if (response.ok) {
    fetchWorkouts(); // Refresh the workout list
  } else {
    const errorData = await response.json();
    console.error("Error creating workout:", errorData);
  }
});

// Add a new set input
document.querySelector("#add-set").addEventListener("click", () => {
  const setsContainer = document.querySelector("#sets-container");
  const setDiv = document.createElement("div");
  setDiv.classList.add("set");
  setDiv.innerHTML = `
    <span class="set-field">
      <label>Weight:</label>
      <input type="number" class="weight" required>
    </span>
    <span class="set-field">
      <label>Reps:</label>
      <input type="number" class="reps" required>
    </span>
  `;
  setsContainer.appendChild(setDiv);
});

// Initialize
fetchWorkouts();
