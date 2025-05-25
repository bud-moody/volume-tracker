const API_URL = "http://localhost:8080/workouts";

// Fetch and display workouts
async function fetchWorkouts() {
  const response = await fetch(API_URL);
  const workouts = await response.json();
  const tableBody = document.querySelector("#workouts-table tbody");
  tableBody.innerHTML = ""; // Clear existing rows

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
  });
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
