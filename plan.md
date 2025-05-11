Goal:
* Create a web app, with persistent storage (postgres db or sqlite?) that allows a user to enter a workout with the following parameters:
** date
** exercise
** a series of sets, and for each set: weight, reps
* The web app should provide visualisations of workout volume over time (i.e. for different dates) where volume is defined as weight x reps. For instance a workout of 50kg for 2 sets of 5 and 60kg for 3 reps would have volume equal to.. 50 * 2 * 5 + 60 * 3 = 500 + 180 = 680
* For the backend please use Go, and for the front end use plain javascript and only one third party library if needed.

Implementation Plan:

1. **Backend Development (Go):**
   - Set up a Go project structure.
   - Use the standard library only for http routing
   - Create a RESTful API to handle CRUD operations for workouts.
     - Endpoints:
       - `POST /workouts`: Add a new workout.
       - `GET /workouts`: Retrieve all workouts.
       - `GET /workouts/{id}`: Retrieve a specific workout.
       - `PUT /workouts/{id}`: Update a workout.
       - `DELETE /workouts/{id}`: Delete a workout.
   - Define models for the database:
     - Workout: `id`, `date`, `exercise`, `sets` (array of sets).
     - Set: `weight`, `reps`.
   - Use a database (SQLite) for persistent storage.
     - Implement database migrations and schema setup.
   - Write unit tests for the API endpoints.

2. **Frontend Development (Plain JavaScript):**
   - Create a simple HTML/CSS layout for the app.
     - Input form for adding workouts.
     - Table or list to display workouts.
     - Section for visualizations.
   - Use JavaScript to interact with the backend API.
     - Fetch and display workouts.
     - Submit new workouts via the form.
     - Update and delete workouts.
   - Use a third-party library (if needed) for data visualization (e.g., Chart.js).
     - Plot workout volume over time.

3. **Integration:**
   - Connect the frontend to the backend.
   - Ensure proper error handling and validation on both frontend and backend.

4. **Deployment:**
   - Containerize the application using Docker.
   - Deploy the app to a cloud platform (e.g., AWS, Heroku, or DigitalOcean).
   - Set up environment variables for database credentials and other configurations.

5. **Testing and Optimization:**
   - Perform end-to-end testing to ensure all features work as expected.
   - Optimize the app for performance and usability.

6. **Documentation:**
   - Write a README file with setup instructions.
   - Document API endpoints and usage.

7. **Future Enhancements:**
   - Add user authentication to allow multiple users to track their workouts.
   - Implement additional visualizations (e.g., progress charts for specific exercises).
   - Allow users to export their workout data.