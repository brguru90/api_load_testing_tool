const chart_option = {
    theme: {
        // monochrome: {
        //   enabled: true
        // }
    },
    plotOptions: {
        pie: {
            dataLabels: {
                offset: -5
            }
        }
    },
    legend: {
        offsetX: 20,
        formatter: function(val, opts) {
            return val + " - " + opts.w.globals.series[opts.seriesIndex]+"%"
        }
    }
};

export {
    chart_option
}