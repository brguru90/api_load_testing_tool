import React from 'react'
import { useSelector } from 'react-redux'

export default function useExtractIteration({ APIindex }) {
    return useSelector(state => {
        const iteration_data = state.metrics_data?.[APIindex]?.iteration_data
        if (iteration_data?.length) {
            return iteration_data
        }
        return []
    })
}
