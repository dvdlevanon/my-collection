import React from 'react'
import { useEffect, useState } from 'react';
import { useParams } from 'react-router-dom'
import '../styles/Pages.css';
import '../styles/ItemPage.css';

function ItemPage() {
    const { itemId } = useParams()
    let [item, setItem] = useState(null);

    useEffect(() => {
		fetch('http://localhost:8080/items/' + itemId)
			.then((response) => response.json())
			.then((item) => setItem(item));
	}, []);

    return (
        <div className='all'>
            <div className='top'>
                <div className="player">
                    { item ? 
                        <video muted controls width="100%" height="700px">
                            <source src={"http://localhost:8080/stream/" + item.url} />
                        </video>
                        : ""
                    }
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
