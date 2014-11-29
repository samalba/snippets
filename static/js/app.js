var postsApp = angular.module('postsApp', []);

postsApp.controller('NavbarCtrl', function($scope, $http) {
    $http.get('/api/user').success(function(data) {
        $scope.user = data;
    });
});
