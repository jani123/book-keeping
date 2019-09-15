require("expose-loader?$!expose-loader?jQuery!jquery");
require("bootstrap/dist/js/bootstrap.bundle.js");
require("bootstrap-datepicker/dist/js/bootstrap-datepicker.js");

$(() => {

    // Update order rows remove button active state
    // Always one order row is left intact
    updateOrderRowRemoveState = function () {
        var $orderRows = $("#order-rows").first();
        var $items = $orderRows.children();
        var $item = null;
        var $button = null;

        $items.each(function (index, value) {
            $item = $(value)
            $button = $item.find("button.order-row-remove").first()
            if (index == 0 && $items.length <= 1) {
                $button.attr("disabled", "true");
                $button.attr("aria-disabled", "true")
            } else {
                $button.attr("disabled", false);
                $button.attr("aria-disabled", false)
            }
        });
    };

    // Add new order row for order
    $("#order-row-add").on("click", function() {
        var $orderRows = $("#order-rows").first();
        var $item = $orderRows.children().first();
        $orderRows.prepend($item.clone(true));
        updateOrderRowRemoveState()
    });

    // Remove order row when remove button is clicked
    $("button.order-row-remove").on("click", function() {
        var $this = $(this);
        $this.parents(".card").remove();
        updateOrderRowRemoveState()
    });

    // Order-rows is loaded, we can do stuff
    $("#order-rows").ready(function () {
        updateOrderRowRemoveState()
    });

    $('input.datepicker').datepicker({
        format: 'dd.mm.yyyy',
        todayBtn: true,
        calendarWeeks: true,
        autoclose: true
    });
});
