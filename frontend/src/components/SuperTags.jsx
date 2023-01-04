import SuperTag from './SuperTag';

function SuperTags({superTags, onSuperTagSelected, onSuperTagDeselected}) {
  return (
    <div className="super-tags">
      {superTags.map((tag) => {
        return <SuperTag key={tag.id} superTag={tag} onSuperTagSelected={onSuperTagSelected} onSuperTagDeselected={onSuperTagDeselected}/>
      })}
    </div>
  )
}

export default SuperTags