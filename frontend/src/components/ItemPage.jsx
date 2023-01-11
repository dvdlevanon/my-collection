import styles from './ItemPage.module.css';
import { useEffect, useState } from 'react';
import { json, useParams } from 'react-router-dom'
import Player from './Player';
import ItemTags from './ItemTags';
import TagChooser from './TagChooser';
import { Dialog } from '@mui/material';

function ItemPage() {
    const { itemId } = useParams()
    let [item, setItem] = useState(null);
    let [addTagMode, setAddTagMode] = useState(false)

    useEffect(() => {
		fetch('http://localhost:8080/items/' + itemId)
			.then((response) => response.json())
			.then((item) => setItem(item));
	}, []);

    const onAddTag = () => {
        setAddTagMode(true);
    }

    const onTagAdded = (tag) => {
        setAddTagMode(false);
        
        item.tags.push(tag)
        fetch('http://localhost:8080/items/' + itemId, {
            method: "POST",
            body: JSON.stringify(item)
        });
    }

    return (
        <div className={styles.all}>
            <div className={styles.top}>
                { item ? <Player item={item} /> : "" }
                { item ? <ItemTags item={item} onAddTag={onAddTag} /> : "" }
            </div>
            <div className={styles.related_items}>
                related-items
            </div>
            <Dialog open={addTagMode} item={item} fullWidth maxWidth="false">
                <TagChooser onTagAdded={onTagAdded} />
            </Dialog>
        </div>
    )
}

export default ItemPage
