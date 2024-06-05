@example
Feature: calculator

# comment

@positive
Scenario: Test addition (2 + 2)
- Add 2 and 2
- Check that result is 4

@negative
Scenario: Test addition, but result is not correct
- Add 2 and 2
- Check that result is 5


@positive
Scenario: Test addition (4 + 4)
- Add 4 and 4
- Check that result is 8


