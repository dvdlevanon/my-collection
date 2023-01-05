import React from 'react'
import ActiveTags from './ActiveTags'
import { CssBaseline, Switch } from '@mui/material';

function SidePanel({activeTags, onTagDeactivated, onTagSelected, onTagDeselected, onChangeCondition }) {

    const onConditionChanged = (e) => {
        onChangeCondition(e.target.checked ? "&&" : "||")
    }

    return (
        <div className="side-panel" >
            { activeTags.length > 1 ? 
            <div className="condition-switch">
                <span>||</span>
                <Switch onChange={(e) => onConditionChanged(e)} />
                <span>&&</span>
            </div>
                : ''
            }
            { activeTags.length > 0 ? 
                <ActiveTags activeTags={activeTags} onTagDeactivated={onTagDeactivated} 
                    onTagSelected={onTagSelected} onTagDeselected={onTagDeselected} />
                : ''
            }
        </div>
    )
}

export default SidePanel