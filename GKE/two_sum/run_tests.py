import os
import json
import requests
import solution

def run_tests(test_cases, user_code):
    # Run your test cases against the user_code here
    # You should return a JSON serializable object with the results
    pass

if __name__ == "__main__":
    user_code = os.environ.get("USER_CODE")
    test_cases_path = "/app/test_cases.json"

    with open(test_cases_path, "r") as f:
        test_cases = json.load(f)

    results = run_tests(test_cases, user_code)

    results_file = "/app/results.json"
    with open(results_file, "w") as f:
        json.dump(results, f)

    # Send the results back to your application
    post_url = "https://api.irvyn.xyz/results-endpoint"
    data = {
        "submission_id": os.environ.get("SUBMISSION_ID"),
        "results": results
    }
    headers = {"Content-Type": "application/json"}
    requests.post(post_url, data=json.dumps(data), headers=headers)
