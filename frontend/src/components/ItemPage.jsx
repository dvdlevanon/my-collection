import styles from './ItemPage.module.css';
import { useEffect, useState } from 'react';
import { useParams } from 'react-router-dom'
import Player from './Player';
import ItemTags from './ItemTags';

function ItemPage() {
    const { itemId } = useParams()
    let [item, setItem] = useState(null);

    useEffect(() => {
		fetch('http://localhost:8080/items/' + itemId)
			.then((response) => response.json())
			.then((item) => setItem(item));
	}, []);

    return (
        <div className={styles.all}>
            <div className={styles.top}>
                { item ? <Player item={item} /> : "" }
                { item ? <ItemTags item={item} /> : "" }
            </div>
            <div className={styles.related_items}>
                related-items
            </div>
        </div>
    )
}

export default ItemPage
