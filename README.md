# Recruitment Management System API

Project Overview: Develop a backend server for a Recruitment Management System where users can create profiles, upload resumes (PDF and DOCX), and apply for jobs. Admins can manage job postings and view applicant data. Uploaded resumes are processed through a third-party API to extract and store relevant information.

## Endpoints

### User Routes

- **POST /signup**
  - **Description:** Create a new user or return an existing user's information.
  - **Request Body:** User details (e.g., name, email, password).

- **POST /login**
  - **Description:** Authenticate user and create a JWT token.
  - **Request Body:** User credentials (e.g., email, password).

- **GET /jobs**
  - **Description:** Retrieve a list of all available jobs.
  - **Response:** List of jobs with details (e.g., job title, description).

- **POST /jobs/apply/:job_id**
  - **Description:** Apply to a specific job by job ID.
  - **Path Parameter:** `job_id` - ID of the job to apply for.
  - **Request Body:** User application details.

### Resume Upload and Retrieval

- **POST /upload-resume**
  - **Description:** Upload a resume to the database as a binary BLOB.
  - **Request Body:** Resume file (e.g., PDF or DOCX).

### Admin Routes

- **POST /admin/job**
  - **Description:** Create a new job posting.
  - **Request Body:** Job details (e.g., title, description).

- **GET /admin/applicants**
  - **Description:** Retrieve a list of all applicants.
  - **Response:** List of applicants with details (e.g., name, email).

- **GET /admin/applicant/:applicant_id**
  - **Description:** Fetch detailed information of a single applicant by applicant ID.
  - **Path Parameter:** `applicant_id` - ID of the applicant.

- **GET /admin/job/:job_id**
  - **Description:** Fetch job and applied applicants details.
  - **Path Parameter:** `job_id` - ID of the job to fetch details for.
  - **Response:** Job details along with a list of applicants who have applied.

### Extra Endpoints ###

- **GET /read-resume**
  - **Description:** Convert binary BLOB to image for the applicant user to view their resume.
  - **Request Body:** None (assumes user is authenticated and has access to their own resume).

- **GET /admin/read-resume**
  - **Description:** For admin, fetch a resume by passing email or user_id in the request body.
  - **Request Body:** JSON object with `email` or `user_id`.

## Installation

**Note**: This project uses MySQL as DBMS. You need to create the necessary tables before running the application. The **create_table_queries.sql** file is provided. Please execute this file using a MySQL query runner or editor to set up the required tables.

1. **Clone the Repository**
   ```bash
   git clone https://github.com/anoopyadav7929/golang-recruitment-system
   ```

2. **Run the Project**
    ```bash
    cd golang-recruitment-system
    go mod tidy 
    go run main.go 
    ```