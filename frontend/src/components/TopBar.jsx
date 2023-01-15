import styles from './TopBar.module.css';
import { Button } from '@mui/material'
import { Link } from 'react-router-dom';

function TopBar() {
    const refreshClicked = (e) => {
        fetch('http://localhost:8080/items/refresh-preview')
    }

    return (
        <div className={styles.top_bar}>
            <Link to="/">
                <Button variant="outlined">Home</Button>
            </Link>
            <Button variant="outlined" onClick={(e) => refreshClicked(e)}>Refresh Gallery</Button>
        </div>
    )
}

export default TopBar