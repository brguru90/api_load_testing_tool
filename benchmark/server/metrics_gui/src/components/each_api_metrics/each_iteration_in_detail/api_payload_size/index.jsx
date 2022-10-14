import React, { useEffect } from 'react'

export default function APIPayLoadSize({APIindex}) {

    useEffect(() => {
        console.log(`Rendered: APIPayLoadSize index=${APIindex}`)
    })


    return (
        <div>
            <h1>APIPayLoadSize</h1>
            <p>Total_request_payload_size_in_bytes</p>
            <p>Total_response_payload_size_in_bytes</p>
        </div>
    )
}
