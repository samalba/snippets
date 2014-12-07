var app = angular.module('snippetsApp', ['ngSanitize', 'ui.router', 'ui.bootstrap', 'ui.ace'])
    .run(['$rootScope', '$state', '$stateParams', '$cacheFactory',
            function($rootScope, $state, $stateParams, $cacheFactory) {
                // Make the state accessible from the root scope
                $rootScope.$state = $state;
                $rootScope.$stateParams = $stateParams;
                $rootScope.cache = $cacheFactory('snippetsCache');
            }
        ]);

app.config(['$stateProvider', '$urlRouterProvider',
    function($stateProvider, $urlRouterProvider) {

        $urlRouterProvider
            .otherwise('/');

        var views = {
            'navbar': {
                templateUrl: 'html/navbar.html',
                controller: 'NavbarCtrl'
            },
            'menu': {
                templateUrl: 'html/menu.html',
                controller: 'MenuCtrl'
            }
        };

        $stateProvider

            .state('index', {
                url: '/',
                views: angular.extend({}, views, {'main': {
                    templateUrl: 'html/index.html'
                }})})

            .state('snippets-new', {
                url: '/snippets/new',
                views: angular.extend({}, views, {'main': {
                    templateUrl: 'html/snippets.edit.html',
                    controller: 'SnippetEditCtrl'
                }})})

            .state('all-teams', {
                url: '/all-teams',
                views: angular.extend({}, views, {'main': {
                    templateUrl: 'html/all-teams.html'
                }})});
    }]);

app.controller('NavbarCtrl', function($scope, $http) {
    var cache = $scope.$root.cache,
        user = cache.get('user');
    if (user) {
        $scope.user = user;
        return;
    }
    $http.get('/api/users/me').success(function(data) {
        cache.put('user', data);
        $scope.user = data;
    });
});

app.controller('MenuCtrl', function($scope, $location) {
    $scope.state = $scope.$root.$state.current.name;
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
