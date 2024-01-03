# Backend API 

## Introduction
Welcome to my Honors Project's Backend API repository, a comprehensive solution designed to facilitate robust authentication, session management, submission storing, and remote code execution. The project is meant to be a similar platform to leetcode, it allows for all of the basic actions you would expect for a code submission website.

## Technologies Used
- **Go (Golang)**: The backbone of the API, chosen for its simplicity, high performance, and excellent support for concurrency.
- **Python**: Utilized for handling specific backend tasks, capitalizing on its vast libraries and ease of use.
- **Docker**: Ensures a seamless and consistent environment across different systems, enhancing the deployment process.

## Features
- **Authentication**: Secure login/signup processes.
- **Session Management**: Efficient handling of user sessions, providing both security and convenience. 
- **Data Storage**: Robust submission storing mechanism, ensuring data integrity and accessibility using both redis and postgres, arguably two of the *best* industry tools.
- **Remote Code Execution**: Facilitates safe and isolated execution of code, a crucial feature for modern development environments. This is done by relying on GKE - running code marking in kubernetes.

## Project Structure
- `.github/workflows`: CI/CD pipeline ensuring code quality and automated testing.
- `GKE`: Google Kubernetes Engine configurations for scalable deployment and for marking the code submissions
- `database`: Database schemas and connection utilities.
- `models`, `utils`: Core modules for business logic and utility functions.
