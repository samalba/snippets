var app = angular.module('snippetsApp', ['ngSanitize', 'ui.router', 'ui.bootstrap', 'ui.ace'])
    .run(['$rootScope', '$state', '$stateParams', '$cacheFactory', '$http',
            function($rootScope, $state, $stateParams, $cacheFactory, $http) {
                // Make the state accessible from the root scope
                $rootScope.$state = $state;
                $rootScope.$stateParams = $stateParams;
                $rootScope.cache = $cacheFactory('snippetsCache');
                $rootScope.getUser = function(scope) {
                    var cache = $rootScope.cache,
                        user = cache.get('user');
                    if (user) {
                        scope.user = user;
                        return;
                    }
                    $http.get('/api/users/me').success(function(data) {
                        cache.put('user', data);
                        scope.user = data;
                    });
                };
            }
        ]);

app.config(['$stateProvider', '$urlRouterProvider',
    function($stateProvider, $urlRouterProvider) {

        $urlRouterProvider
            .otherwise('/');

        var views = {
            'navbar': {
                templateUrl: 'html/navbar.html',
                controller: 'UserCtrl'
            },
            'menu': {
                templateUrl: 'html/menu.html',
                controller: 'UserCtrl'
            }
        };

        $stateProvider

            .state('index', {
                url: '/',
                views: angular.extend({}, views, {'main': {
                    templateUrl: 'html/index.html'
                }})})

            .state('settings', {
                url: '/users/settings',
                views: angular.extend({}, views, {'main': {
                    templateUrl: 'html/users-settings.html',
                    controller: 'UserCtrl'
                }})})

            .state('snippets-new', {
                url: '/snippets/new',
                views: angular.extend({}, views, {'main': {
                    templateUrl: 'html/snippets-edit.html',
                    controller: 'SnippetEditCtrl'
                }})})

            .state('admin-users', {
                url: '/admin/users/:login',
                views: angular.extend({}, views, {'main': {
                    templateUrl: 'html/admin-users.html',
                    controller: 'AdminUsersCtrl'
                }})})

            .state('admin-orgs', {
                url: '/admin/orgs',
                views: angular.extend({}, views, {'main': {
                    templateUrl: 'html/admin-orgs.html'
                }})})

            .state('admin-teams', {
                url: '/admin/teams',
                views: angular.extend({}, views, {'main': {
                    templateUrl: 'html/admin-teams.html'
                }})})

            .state('all-teams', {
                url: '/all-teams',
                views: angular.extend({}, views, {'main': {
                    templateUrl: 'html/all-teams.html'
                }})});
    }]);

app.controller('UserCtrl', function($scope) {
    $scope.state = $scope.$state.current.name;
    $scope.getUser($scope);
});

app.controller('AdminUsersCtrl', function($scope, $http) {
    $http.get("/api/users").success(function(data) {
        $scope.users = data;
    });

    console.log($scope.$stateParams);

    $scope.create = function(user) {
        console.log(user);
    };
});

app.controller('SnippetEditCtrl', function($scope) {

    $scope.preview = "";

    $scope.aceLoaded = function(editor) {
        editor.setTheme("ace/theme/textmate");
        editor.getSession().setMode("ace/mode/markdown");
        editor.focus();
        $scope.editor = editor;
    };

    $scope.renderPreview = function() {
        if (!$scope.editor) {
            return;
        }
        $scope.preview = markdown.toHTML($scope.editor.getValue());
        if (!$scope.preview) {
            $scope.preview = '<span class="italic">no preview available</span>';
        }
    }
});
