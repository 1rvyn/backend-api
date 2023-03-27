import os
import json
import requests
import solution

def run_tests(test_cases, user_code):
    results = []

    # Write user code to a temporary file
    with open("user_code.py", "w") as f:
        f.write(user_code)

    # Import the user's code
    import user_code as user_solution

    # Iterate through test cases
    for i, test_case in enumerate(test_cases):
        # Get the input and expected output
        input_data = test_case["input"]
        expected_output = test_case["output"]

        try:
            # Run the user's function
            user_output = user_solution.two_sum(*input_data)
        except Exception as e:
            results.append({
                "test_case": i + 1,
                "status": "failed",
                "error": str(e)
            })
            continue

        # Compare the user's output with the expected output
        if user_output == expected_output:
            results.append({
                "test_case": i + 1,
                "status": "passed"
            })
        else:
            results.append({
                "test_case": i + 1,
                "status": "failed",
                "expected": expected_output,
                "received": user_output
            })

    return results
