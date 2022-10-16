---
title: "Roadmap"
weight: 2
# bookFlatSection: false
bookToc: true
# bookHidden: false
# bookCollapseSection: false
# bookComments: false
# bookSearchExclude: false
---

# Roadmap

Everything is a work in progress. Feedback and suggestions are very welcome!
Open a [ticket on
GitHub](https://github.com/cleodora-forecasting/cleodora/issues) or [send an
e-mail (info@cleodora.org)](mailto:info@cleodora.org).

## Version 1

The MVP (Minimal Viable Product) must implement the following features.

* Web-based
* Make binary (yes/no) predictions (title, description, probability of being
  true 0 - 100, datetime, reason for probability, resolution date)
* Edit question. Modify prediction title, description and resolution date
* Update prediction. Allow modifying the probability 0 - 100 and attach a
  datetime to it. Add a reason for that probability.
* Allow viewing the history of probabilities for a prediction, including
  datetimes and the reason given at the time.
* Allow resolving a prediction as true, false or N/A. Add a comment when
  resolving.
* Calculate the Brier score for every probability of a prediction in the past
* Calculate the avg Brier score for the entire prediction
* Calculate the avg Brier score across all predictions
* For personal usage at home
* Single binary, easy deployment cross platform desktop/server
* High degree of test automation and good documentation to make future
  development easy for busy and lazy people like myself
* Have an overview page with all predictions.
	* Sortable by:
		* Last prediction update
		* Creation date
		* Last edit
		* Resolution date
		* Brier score
		* Resolution
	* Filter by:
		* Resolution
* Display some statistics, possibly on the home page:
	* Current average Brier score and evolution over time
	* Predictions added over time
	* Number of updates over time


## Possible later features

* Allow choosing the datetime (default: now) when updating the probability of a
  prediction
* Have some kind of categories, possibly hierarchical to see Brier score for
  categories
* Version predictions whenever the title, description or resolution date is
  updated
* Probabilities always refer to a version
* When displaying a past probability display the version of the question at the
  time, in particular the resolution date set at the time (which indicates how
  much uncertainty there was at the time)
* Allow seeing avg Brier score X time before the resolution date (as known at
  the time) i.e. what is my avg Brier score two months before resolution, even
  if for one of the questions the resolution came 1 week later initially the
  resolution was predicted later.
* Multiple users for them to compare scores
* Allow having groups of people e.g. for company prediction tournaments
* Public instance where people can make predictions for other people to see
	* app.cleodora.org ?
    * Allow marking predictions as public, private or private within a group of
      people
* Always maintain self-deployment option for making personal predictions you
  don't want to share with the Internet (e.g. personal relationships, personal
  goals, etc.)
* More scoring systems, not only Brier score
* Data import and export to transfers between private and public instances and
  to analyze with other tools
* Allow exporting/importing predictions as JSON to share with other people so
  they can predict in their own Cleodora instance. What to do when editing a
  question? New export that contains version information and can be merged?
* Allow exporting nice reports or graphs to brag and show other people
* Multi-answer questions. Always have one answer 'other' that has 100 -
  sum(others) by default.
    * When editing a question and adding new answers then those new answers get
      0 probability initially. This can lead to drastic changes in Brier score
      (if 'other' was rather large before), but that is fine.
