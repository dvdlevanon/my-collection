import { Button } from '@mui/material'
import React from 'react'

function TopBar() {
    const refreshClicked = (e) => {
        fetch('http://localhost:8080/items/refresh-preview')
    }

    return (
        <div className='top-bar'>
            <Button variant="outlined" onClick={(e) => refreshClicked(e)}>Refresh Gallery</Button>
        </div>
    )
}

export default TopBar