---
title: "Use Cases"
weight: 1
# bookFlatSection: false
# bookToc: true
# bookHidden: false
# bookCollapseSection: false
# bookComments: false
# bookSearchExclude: false
---

# Use Cases

Everything is a work in progress. Feedback and suggestions are very welcome!
Open a [ticket on
GitHub](https://github.com/cleodora-forecasting/cleodora/issues) or [send an
e-mail (info@cleodora.org)](mailto:info@cleodora.org).

Some examples how Cleodora could be used.

## Mindy the engineer

Mindy is an engineer who works at a company that just hired a new CEO on
January 1st 2022. She predicts he will leave the company by the end of the
year.

Mindy makes the following forecast on January 2nd 2022:

* Will Tom, our new CEO, leave the company by December 31st 2022?
* Prediction: Yes
* Confidence: 60%
* Reasoning: His track record shows that he switched companies 8 times in the
  last 10 years and I have a bad gut feeling about him.

On March 12th 2022 Mindy updates her forecast:

* Confidence: 60% → 70%
* Reasoning: In yesterday's meeting Tom was unable to name our three top
  customers. This shows he doesn't care about the company.

On June 30th 2022 Mindy updates her forecast again:

* Confidence: 70% → 30%
* Reasoning: He has been showing much more interest in ongoing projects and
  started engaging with team leads in strategic decision making.

On August 1st 2022 Mindy updates her forecast again:

* Confidence: 30% → 85%
* Reasoning: I am sure he will leave, just not completely sure whether it will
  be before the end of the year. The door to this office was not closed and I
  overheard him talking to a competitor.

On November 4th 2022 Mindy resolves the forecast to "Yes":

* Reasoning: Tom just announced that he will leave by November 30th.

Cleodora now calculates Mindy's Brier score for her forecast. The Brier score
is a common measure for forecast accuracy where 0 is best and 1 is worst:

* January 2nd 2022 (11 months, 29 days before resolution): 0.16
* March 12th 2022 (9 months, 19 days before resolution): 0.09
* June 30th 2022 (6 months, 1 day before resolution): 0.49
* August 1st 2022 (4 months, 30 days before resolution): 0.02
* Overall average: 0.19

Mindy analyzes how her Brier score evolved over time and sees that it tanked on
June 30th. Reading her reasoning at the time she decides that she should not
have adjusted down her prediction so strongly given the other evidence she had.
Instead she should have adjusted to maybe 55%. In the future she will be more
careful in such situations. In general she will trust her gut feeling about
people a little more, because she was right in her forecast after all.
