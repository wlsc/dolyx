$(document).ready(function () {

    $(".remove-all-images").click(function () {
        sendControlCommand("images", {"name": "removeall", "value": "all"});
    });

    $('button[data-target="#removalModal"]').click(function () {

        var imageId = $(this).data("id");
        $("#dockerImageTag").text($(this).data("tag"));

        $(".remove-image").on("click", function () {
            sendControlCommand("images", {"name": "remove", "value": imageId});
        });
    });

    $('#pruneContainers').click(function () {
        sendControlCommand("prune", {"name": "containers", "value": null});
    });

    $('#pruneImages').click(function () {
        sendControlCommand("prune", {"name": "images", "value": null});
    });

    $('#pruneVolumes').click(function () {
        sendControlCommand("prune", {"name": "volumes", "value": null});
    });

    $('#pruneNetworks').click(function () {
        sendControlCommand("prune", {"name": "networks", "value": null});
    });

    $('#pruneCache').click(function () {
        sendControlCommand("prune", {"name": "cache", "value": null});
    });

    $('#pruneAll').click(function () {
        sendControlCommand("prune", {"name": "all", "value": null});
    });

    function sendControlCommand(commandType, localCommand) {

        $('#workingLabel').show();

        $.ajax({
            method: "POST",
            url: "/control",
            data: JSON.stringify({
                type: commandType,
                command: {
                    name: localCommand.name,
                    value: localCommand.value
                }
            }),
            contentType: "application/json",
            dataType: 'json'
        }).done(function () {
            location.reload();
        }).fail(function (response) {
            $('#workingLabel').hide();
            handleServerErrors(response);
        });
    }

    function handleServerErrors(response) {

        console.log(JSON.stringify(response));
        var response = $.parseJSON(response.responseText);
        console.log(response.error);

        var status = response.status;

        switch (status) {
            case 409:
                alert(response.error);
                break;
        }
    }
});