var myApp = angular.module('myApp', []);

function autoCompleteViewController($scope, $http) {
	$scope.words = ["test word1", "word2"];
	$scope.suggestions = ["test suggestion1", "test suggestion2"];
}

function searchResultsViewController($scope, $http) {
	$scope.results = ["result1", "result2"];
}

myApp.controller('autoCompleteViewController', autoCompleteViewController);
myApp.controller('searchResultsViewController', searchResultsViewController);
