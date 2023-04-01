#!/usr/bin/env python3
import os
import unittest
from io import StringIO
from unittest import TextTestRunner
from flask import Flask, request, jsonify

app = Flask(__name__)

@app.route('/run_tests', methods=['POST'])
def run_tests():
    if 'code' not in request.files:
        return jsonify({"error": "Missing 'code' file in the request"}), 400
    code = request.files['code']
    code.save('submitted_code.py')

    # Run the tests
    suite = unittest.TestLoader().loadTestsFromName('test_two_sum')
    with StringIO() as test_output:
        test_result = TextTestRunner(stream=test_output, verbosity=2).run(suite)

    # Parse test results
    test_results = {
        f"test {i + 1}": not any(case == failed_case for failed_case, _ in test_result.failures) and
                         not any(case == error_case for error_case, _ in test_result.errors)
        for i, case in enumerate(suite)
    }

    # Delete submitted code file
    os.remove('submitted_code.py')

    return jsonify({"result": test_results})



if __name__ == '__main__':
    app.run(host='0.0.0.0', port=8080)
