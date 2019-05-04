{{ template "header.tpl" .}}
<div class="container">
    <div class="row">
        <table class="table table-striped table-hover col-xs-offset-4">
            <thead>
            <tr>
                <th>#</th>
                <th>Tag</th>
                <th>Id</th>
                <th>Created</th>
                <th>Size</th>
                <th>Control</th>
            </tr>
            </thead>
            <tbody>
            {{ range $index, $image := .images }}
            <tr>
                <td>{{ $index }}</td>
                <td>{{ $image.Tag }}</td>
                <td title="{{ $image.Id }}">{{ printf "%.12s" $image.Id }}</td>
                <td title="{{ $image.Created }}">{{ $image.Created }}</td>
                <td>{{ $image.Size }}</td>
                <td>
                    <button class="btn btn-danger btn-sm" title="Remove image" data-id="{{$image.Id}}" data-tag="{{$image.Tag}}" data-toggle="modal" data-target="#removalModal">
                        <i class="material-icons">remove_circle_outline</i>
                    </button>
                </td>
            </tr>
            {{ end }}
            </tbody>
        </table>
        {{if not .images}}
        <div class="col-xs-12 offset-5">
            No images found.
        </div>
        {{end}}
        {{if .images}}
        <div class="col-xs-12">
            <button class="btn btn-danger btn-sm" style="float: right;" title="Stop containers and remove all images" data-toggle="modal" data-target="#removalAllModal">
                Remove all images
            </button>
        </div>
        {{end}}
    </div>
</div>
{{if .images}}
<div class="modal fade" id="removalModal" tabindex="-1" role="dialog" aria-labelledby="removalModalTitle" aria-hidden="true">
    <div class="modal-dialog modal-dialog-centered" role="document">
        <div class="modal-content">
            <div class="modal-header">
                <h5 class="modal-title" id="removalModalTitle"><i class="material-icons">remove_circle_outline</i> Removal of Docker image</h5>
                <button type="button" class="close" data-dismiss="modal" aria-label="Close">
                    <span aria-hidden="true">&times;</span>
                </button>
            </div>
            <div class="modal-body">
                Are you <strong>really</strong> want to remove <strong id="dockerImageTag">unknown image name</strong>?
            </div>
            <div class="modal-footer">
                <button type="button" class="btn btn-secondary" data-dismiss="modal">Close</button>
                <button type="button" class="btn btn-danger remove-image">Remove</button>
            </div>
        </div>
    </div>
</div>
<div class="modal fade" id="removalAllModal" tabindex="-1" role="dialog" aria-labelledby="removalAllModalTitle" aria-hidden="true">
    <div class="modal-dialog modal-dialog-centered" role="document">
        <div class="modal-content">
            <div class="modal-header">
                <h5 class="modal-title" id="removalAllModalTitle"><i class="material-icons">remove_circle_outline</i> Removal of Docker images</h5>
                <button type="button" class="close" data-dismiss="modal" aria-label="Close">
                    <span aria-hidden="true">&times;</span>
                </button>
            </div>
            <div class="modal-body">
                Are you <strong>really</strong> want to remove <strong>all</strong> images?
            </div>
            <div class="modal-footer">
                <button type="button" class="btn btn-secondary" data-dismiss="modal">Close</button>
                <button type="button" class="btn btn-danger remove-all-images">Remove all</button>
            </div>
        </div>
    </div>
</div>
{{end}}
{{ template "footer.tpl" .}}