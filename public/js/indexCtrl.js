'use strict'
var myApp = angular.module('myApp', []);

myApp.controller('IndexCtrl', ['$scope', '$interval', '$timeout', function ($scope, $interval, $timeout) {
    var quotes = ["640K ought to be enough for anybody", "The architects are still drafting", "The bits are breeding", "We're building the buildings as fast as we can", "Would you prefer chicken, steak, or tofu?", "Pay no attention to the man behind the curtain", "Enjoy the elevator music", "While the little elves draw your map", "A few bits tried to escape, but we caught them", "And dream of faster computers", "Would you like fries with that?", "Checking the gravitational constant in your locale", "Go ahead -- hold your breath", "At least you're not on hold", "Hum something loud while others stare", "You're not in Kansas any more", "The server is powered by a lemon and two electrodes", "We love you just the way you are", "While a larger software vendor in Seattle takes over the world", "We're testing your patience", "As if you had any other choice", "Take a moment to sign up for our lovely prizes", "Don't think of purple hippos", "Follow the white rabbit", "Why don't you order a sandwich?", "While the satellite moves into position", "The bits are flowing slowly today", "Dig on the 'X' for buried treasure... ARRR!", "It's still faster than you could draw it"];
    $scope.pageIndex = 1;
    $scope.loading = false;
    $scope.quote = quotes[Math.floor(Math.random() * quotes.length)];

    $scope.data = {
        files: [],
        commits: [],
        stats: []
    };

    var quoteInterval;

    $scope.search = function () {
        $scope.loading = true;

        quoteInterval = $interval(function () {
            $scope.quote = quotes[Math.floor(Math.random() * quotes.length)];
        }, 2000);

        $timeout(function () {
            $interval.cancel(quoteInterval);
            $scope.loading = false;
            $scope.pageIndex = 2;

            $scope.data.files = [
                "This is a file name1",
                "This is a file name2",
                "This is a file name3",
                "This is a file name4",
                "This is a file name5",
                "This is a file name6",
                "This is a file name7",
                "This is a file name8"
            ];

            $scope.data.commits = [
                "This is a commit name1",
                "This is a commit name2",
                "This is a commit name3",
                "This is a commit name4",
                "This is a commit name5",
                "This is a commit name6",
                "This is a commit name7",
                "This is a commit name8"
            ];

            $scope.data.stats = [
                "This is a stat name1",
                "This is a stat name2",
                "This is a stat name3",
                "This is a stat name4",
                "This is a stat name5",
                "This is a stat name6",
                "This is a stat name7",
                "This is a stat name8"
            ];
        }, 4000);
    };

}]);