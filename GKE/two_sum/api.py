import os
import sys
import io
import tempfile
import importlib.util
from contextlib import redirect_stdout
from flask import Flask, request, jsonify
import unittest

app = Flask(__name__)

class TestTwoSum(unittest.TestCase):
    def test_case_1(self):
        nums = [2, 7, 11, 15]
        target = 9
        output = [0, 1]
        self.assertEqual(two_sum(nums, target), output)

    def test_case_2(self):
        nums = [3, 2, 4]
        target = 6
        output = [1, 2]
        self.assertEqual(two_sum(nums, target), output)

    def test_case_3(self):
        nums = [3, 3]
        target = 6
        output = [0, 1]
        self.assertEqual(two_sum(nums, target), output)

@app.route('/run_tests', methods=['POST'])
def run_tests():
    code = request.files['code']

    with tempfile.NamedTemporaryFile(suffix=".py", delete=False) as temp:
        temp.write(code.stream.read())
        temp.flush()

        spec = importlib.util.spec_from_file_location("submitted_code", temp.name)
        submitted_code = importlib.util.module_from_spec(spec)
        spec.loader.exec_module(submitted_code)
        global two_sum
        two_sum = submitted_code.two_sum

    suite = unittest.TestLoader().loadTestsFromTestCase(TestTwoSum)
    test_results = []

    for case in suite._tests:
        method = case._testMethodName
        with io.StringIO() as buffer:
            with redirect_stdout(buffer):
                result = case.run()
                test_output = buffer.getvalue()

        test_results.append({
            "test_name": method,
            "success": result.wasSuccessful(),
            "output": test_output.strip()
        })

    os.unlink(temp.name)  # Delete the temporary file
    print(test_results)
    return jsonify(test_results)


if __name__ == '__main__':
    app.run(host='0.0.0.0', port=8080)
