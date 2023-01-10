import React from 'react'
import { useParams } from 'react-router-dom'
import '../styles/Pages.css';

function ItemPage({  }) {
    const { itemId } = useParams()

    return (
        <div className='all'>
            <div className='top'>
                <div className="player">
                    player
                </div>
                <div className='tags-editor'>
                    tag-editor
                </div>
            </div>
            <div className='related-items'>
                related-items
            </div>
        </div>
    )
}

export default ItemPage