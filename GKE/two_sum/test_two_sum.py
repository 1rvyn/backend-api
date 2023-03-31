import unittest
from submitted_code import two_sum

class TestTwoSum(unittest.TestCase):
    def test_cases(self):
        test_cases = [
            {
                "input": {
                    "nums": [2, 7, 11, 15],
                    "target": 9
                },
                "output": [0, 1]
            },
            {
                "input": {
                    "nums": [3, 2, 4],
                    "target": 6
                },
                "output": [1, 2]
            },
            {
                "input": {
                    "nums": [3, 3],
                    "target": 6
                },
                "output": [0, 1]
            }
        ]

        for case in test_cases:
            self.assertEqual(two_sum(case["input"]["nums"], case["input"]["target"]), case["output"])

if __name__ == '__main__':
    unittest.main()
