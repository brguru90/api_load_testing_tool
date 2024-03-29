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
        fontSize: '9px',
        offsetX: 20,
        formatter: function(val, opts) {
            return val + " - " + opts.w.globals.series[opts.seriesIndex]+"%"
        }
    },
    dataLabels:{
        style: {
            fontSize: '9px',
        },
    },
    title: {
        text: ['Average occurrence of status code ','from all iteration'],
        align: 'left',
        floating: false,
    },
    chart: {
        width: '100%'
    }
};

export {
    chart_option
}