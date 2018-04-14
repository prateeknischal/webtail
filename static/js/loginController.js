angular.module("webtail").controller("loginController", loginController);

loginController.$inject = ["$rootScope", "$scope"]

function loginController($rootScope, $scope) {
  var vm = this;

  $scope.showMessage = false
  $scope.success = true;
  function getUrlVars(){
    var vars = [], hash;
    var hashes = window.location.href.slice(window.location.href.indexOf('?') + 1).split('&');
    for(var i = 0; i < hashes.length; i++)
    {
        hash = hashes[i].split('=');
        vars.push(hash[0]);
        vars[hash[0]] = hash[1];
    }
    return vars;
  }
  vm.init = function init() {
    console.log("In the login controller")
    var query = getUrlVars();
    console.log(query);
    if (query["err"] === "invalid") {
      $scope.success = false;
      $scope.message = "Invalid Credentials"
      $scope.text_color="red"
      $scope.showMessage = true;
    }
    else if (query["logout"] === "success") {
      $scope.success = true;
      $scope.message = "Logout Successful"
      $scope.text_color="green"
      $scope.showMessage = true;
    }
  }

  vm.init();
}
