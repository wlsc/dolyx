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
                    <button class="btn btn-danger btn-sm remove-image" title="Remove image" data-id="{{$image.Id}}">
                        <i class="material-icons">remove_circle_outline</i>
                    </button>
                </td>
            </tr>
            {{ end }}
            </tbody>
        </table>
        <div class="col-xs-12">
            <button class="btn btn-danger btn-sm remove-all-images" style="float: right;" title="Remove all images">
                <i class="material-icons">remove_circle_outline</i> Remove all images
            </button>
        </div>
    </div>
</div>
{{ template "footer.tpl" .}}