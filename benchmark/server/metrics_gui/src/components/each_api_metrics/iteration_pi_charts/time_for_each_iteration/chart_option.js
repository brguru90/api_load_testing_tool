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
        offsetX: -20,
        formatter: function(val, opts) {
            return val + " - " + opts.w.globals.series[opts.seriesIndex]+" Sec"
        }
    },
    dataLabels:{
        style: {
            fontSize: '10px',
            fontFamily: 'Helvetica, Arial, sans-serif',
            fontWeight: '100',
            colors: undefined
        },
        formatter: function(val, opts) {
            return  opts.w.globals.series[opts.seriesIndex]+" Sec"
        }
    },
    title: {
        text: 'Time taken to complete all individual iterations',
        align: 'center',       
        floating: false,
    },
};

export {
    chart_option
}