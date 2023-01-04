import React from 'react'
import ActiveTag from './ActiveTag'

function ActiveTags({activeTags, onTagDeactivated, onTagSelected, onTagDeselected }) {
  return (
    <div className='active-tags'>
        {activeTags.map((tag) => {
            return <ActiveTag key={tag.id} tag={tag} onTagDeactivated={onTagDeactivated}
              onTagSelected={onTagSelected} onTagDeselected={onTagDeselected} />
        })}
    </div>
  )
}

export default ActiveTags