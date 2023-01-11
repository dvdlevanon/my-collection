import styles from './TopBar.module.css';
import { Button } from '@mui/material'

function TopBar() {
    const refreshClicked = (e) => {
        fetch('http://localhost:8080/items/refresh-preview')
    }

    return (
        <div className={styles.top_bar}>
            <Button variant="outlined" onClick={(e) => refreshClicked(e)}>Refresh Gallery</Button>
        </div>
    )
}

export default TopBar