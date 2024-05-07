@example
Feature: calculator

# comment

@positive
Scenario: Test addition
- Add 2 and 2
- Check that result is 4

@positive
Scenario: Test addition
- Add 4 and 4
- Check that result is 8

@positive
Scenario: Test addition (strings)
- Add string "asd" and "fgh"
- Check that result string is "asdfgh"


@negative
Scenario: Test addition, but result is not correct
- Add 2 and 2
- Check that result is 5

