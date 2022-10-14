const chart_option = {
    labels: ["Connecting", "Processing"],
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
        fontSize: '9px',
        offsetX: 20,
        // formatter: function(val, opts) {
        //     return val + " - " + opts.w.globals.series[opts.seriesIndex]+" ms"
        // }
    },
    dataLabels: {
        style: {
            fontSize: '9px',
        },
        formatter: function (val, opts) {
            return opts.w.globals.series[opts.seriesIndex] + " ms"
        }
    },
    title: {
        text: 'Average of API timings for all Iterations',
        align: 'left',        
        floating: false,
    },
};

export {
    chart_option
}