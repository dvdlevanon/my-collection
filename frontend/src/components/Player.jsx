import React from 'react'

function Player({item}) {
    return (
        <div className="player">
            <video muted controls width="100%" height="400px">
                <source src={"http://localhost:8080/stream/" + item.url} />
            </video>
        </div>
    )
}

export default Player