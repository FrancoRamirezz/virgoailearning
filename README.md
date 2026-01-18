# VirgoAI Learning

VirgoAI Learning is a web-based learning platform designed to support immigrants on their journey toward U.S. citizenship and long-term integration. The application focuses on structured citizenship test preparation, English practice for real interviews, and readiness tracking to help users feel confident and prepared.

This repository is organized as a **monorepo**, containing both the frontend (Next.js) and backend (Go) applications.



## Purpose

Many immigrants preparing for U.S. citizenship feel overwhelmed, unsure of what to study, and uncertain about whether they are ready for the citizenship interview.

VirgoAI Learning aims to:

* Provide clear, structured citizenship study materials
* Offer English practice focused on real interview scenarios
* Track learning progress and readiness over time
* Reduce anxiety by showing users exactly where they stand and what to study next



## Who This App Is For

* Immigrants preparing for the U.S. citizenship test
* Adult learners studying civics and interview English
* Families and community organizations supporting citizenship preparation
* Educators or programs looking for structured readiness tools



## Current Features (MVP)

Frontend:

* Authentication pages (login and signup UI)
* Landing and informational pages
* Course and dashboard UI structure
* Responsive design with a focus on accessibility
* Readiness-based messaging and onboarding flow

Backend (in progress):

* Go-based API service
* Authentication and session handling (planned)
* Integration with a PostgreSQL database (planned)
* Support for civics content, quizzes, and progress tracking (planned)

The backend is actively under development and will evolve as the project progresses.



## Repository Structure

```
.
├── frontend/    # Next.js frontend application
├── backend/     # Go backend application
├── README.md    # Project documentation
```



## Getting Started

### Prerequisites

* Node.js (v18 or later recommended)
* npm
* Go (version will be finalized by the backend team)



## Running the Frontend

Navigate to the frontend directory and start the development server:

```bash
cd frontend
npm install
npm run dev
```

The frontend will be available at:

```
http://localhost:3000
```



## Running the Backend

The backend is a Go service and is currently under active development.

Basic placeholder workflow (will be updated as the backend stabilizes):

```bash
cd backend
go run ./cmd/api
```

More detailed instructions, environment variables, and database setup will be added once the backend structure is finalized.



## Development Workflow

* This is a shared team repository.
* All development should be done on feature branches.
* Direct commits to `main` should be avoided.
* Changes should be merged via pull requests.

Recommended branch naming:

* `frontend/<feature-name>`
* `backend/<feature-name>`



## Environment Configuration

Environment variables will be documented as backend and frontend integration progresses.

Frontend environment variables will use:

```
NEXT_PUBLIC_*
```

Backend environment variables will be documented separately.



## Roadmap (High-Level)

* Complete backend authentication and session management
* Implement database schema for courses, quizzes, and progress tracking
* Connect frontend dashboard to real backend data
* Add onboarding quiz and readiness scoring
* Introduce premium features and subscription handling
* Expand platform into a broader immigrant onboarding and learning hub



## License

License information will be added once the project reaches a stable public release.


