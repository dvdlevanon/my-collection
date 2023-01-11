import React from 'react'
import ItemTag from './ItemTag'

function ItemTags({item}) {
  return (
    <div className="item-tags">
        {item.tags.map((tag) => {
            return <ItemTag key={tag.id} tag={tag} />
        })}
    </div>
  )
}

export default ItemTags