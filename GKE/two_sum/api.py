from flask import Flask, request
import requests
import subprocess
import json

app = Flask(__name__)

@app.route('/run_tests', methods=['POST'])
def run_tests():
    code = request.files['code']
    code.save('submitted_code.py')

    result = subprocess.run(["python", "run_tests.py"], capture_output=True, text=True)
    test_output = result.stdout

    # Parse test_output into a dictionary of test results
    test_results = parse_test_output(test_output)

    response = requests.post('https://api.irvyn.xyz/tested', json=test_results)

    return test_results

def parse_test_output(test_output):
    # Implement this function to parse the test_output string
    # and create a dictionary with the format:
    # {
    #   "test 1": True,
    #   "test 2": True,
    #   "test 3": False,
    #   ...
    # }
    pass

if __name__ == '__main__':
    app.run(host='0.0.0.0', port=8080)
