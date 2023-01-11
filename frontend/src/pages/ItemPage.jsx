import React from 'react'
import { useEffect, useState } from 'react';
import { useParams } from 'react-router-dom'
import '../styles/Pages.css';
import '../styles/ItemPage.css';
import Player from '../components/Player';
import ItemTags from '../components/ItemTags';

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
                { item ? <Player item={item} /> : "" }
                { item ? <ItemTags item={item} /> : "" }
            </div>
            <div className='related-items'>
                related-items
            </div>
        </div>
    )
}

export default ItemPage
