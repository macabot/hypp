import React, { useState } from 'react';

const lists = [
  [{id: 1, label: 'A', pos: 0}, {id: 2, label: 'B', pos: 1}, {id: 3, label: 'C', pos: 2}],
  [{id: 2, label: 'B', pos: 0}, {id: 3, label: 'C', pos: 1}, {id: 1, label: 'A', pos: 2}],
]

export function App(props) {
  const [index, setIndex] = useState(0)
  const list = [...lists[index]]
  // list.sort((a, b) => a.id - b.id) // Why is this needed?
  return (
    <div className='App'>
      <button onClick={() => setIndex((index + 1) % lists.length)} >Next</button>
      <div className="wrapper">
        {
          list.map(item => (
            <div key={item.id} className={`item position-${item.pos}`}>{item.label}</div>
          ))
        }
      </div>
    </div>
  );
}
