const API_URL = "http://localhost:8080/workouts";

// Render the chart
function renderChart(data) {
  const ctx = document.getElementById("volume-chart").getContext("2d");

  // Destroy the previous chart instance if it exists
  if (window.volumeChart) {
    window.volumeChart.destroy();
  }

  // Create a new chart
  window.volumeChart = new Chart(ctx, {
    type: "line",
    data: {
      labels: data.labels, // Dates
      datasets: [{
        label: "Workout Volume",
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
    const chartData = {
      labels: [], // Dates
      volumes: [] // Workout volumes
    };

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

      // Add data to chart
      chartData.labels.push(workout.date);
      chartData.volumes.push(volume);
    });

    // Reverse the order of labels and volumes to ensure time moves from left to right
    chartData.labels.reverse();
    chartData.volumes.reverse();

    // Render the chart
    renderChart(chartData);
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
    <label>Weight:</label>
    <input type="number" class="weight" required>
    <label>Reps:</label>
    <input type="number" class="reps" required>
  `;
  setsContainer.appendChild(setDiv);
});

// Initialize
fetchWorkouts();
