$(document).ready(function() {

    $(".remove-all-images").on("click", function() {

        $.ajax({
          method: "POST",
          url: "/control",
          data: JSON.stringify({
                  type: "images",
                  command : {
                         name: "removeall",
                         value: "all"
                         }
                }),
          contentType: "application/json",
          dataType: 'json'
        })
        .done(function(status) {
            location.reload();
          })
        .fail(function(response) {
            handleError(response);
        });
    });

    $(".remove-image").on("click", function() {
        var id = $(this).data("id");

        $.ajax({
          method: "POST",
          url: "/control",
          data: JSON.stringify({
                  type: "images",
                  command : {
                         name: "remove",
                         value: id
                         }
                }),
          contentType: "application/json",
          dataType: 'json'
        })
        .done(function(status) {
            location.reload();
          })
        .fail(function(response) {
            handleError(response);
        });
    });

    /**
     *  Handles HTTP error codes from server's response
     */
    function handleError(response) {

        console.log(JSON.stringify(response));
        var response = $.parseJSON(response.responseText);
        console.log(response.error);

        var status = response.status;

        switch(status){
            case 409:
                alert(response.error);
                break;
        }
    }
});