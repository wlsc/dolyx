{{ template "header.tpl" .}}
<div class="container">
    <div class="row">
        <table class="table table-striped table-hover col-xs-offset-4">
            <thead>
                <tr>
                    <th>#</th>
                    {{ range $header_index, $name := .images_list_header }}
                    <th>{{ $name }}</th>
                    {{ end }}
                    <th>Control</th>
                </tr>
            </thead>
            <tbody>
                {{ range $row_index, $row := .images_list }}
                <tr>
                    <td>{{ $row_index }}</td>
                    {{ range $col_index, $col := $row }}
                    <td>{{ $col }}</td>
                    {{ end }}
                    <td>
                        <button class="btn btn-danger btn-sm remove-image" title="Remove image" data-id="{{index $row 2}}">
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