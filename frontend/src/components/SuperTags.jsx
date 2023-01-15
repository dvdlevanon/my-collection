import styles from './SuperTags.module.css';
import SuperTag from './SuperTag';

function SuperTags({superTags, onSuperTagClicked}) {
    return (
        <div className={styles.super_tags}>
            {superTags.map((tag) => {
                return <SuperTag key={tag.id} superTag={tag} onSuperTagClicked={onSuperTagClicked} />
            })}
        </div>
    )
}

export default SuperTags