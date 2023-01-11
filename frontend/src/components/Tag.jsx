import styles from './Tag.module.css';

function Tag({tag, onTagActivated}) {
    return (
        <>
            <div className={styles.tag + " " + (tag.selected ? styles.selected : styles.unselected)} onClick={() => onTagActivated(tag)}>
                {tag.title}
            </div>
        </>
    )
}

export default Tag
