'use strict';
/*!
 * @license
 * Copyright 2019-2020 The Go Authors. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
const toggles = document.querySelectorAll('[data-toggletip-content]');
toggles.forEach(toggle => {
  const message = toggle.getAttribute('data-toggletip-content');
  const tip = toggle.nextElementSibling;
  toggle.addEventListener('click', () => {
    if (!tip) {
      return;
    }
    tip.innerHTML = '';
    setTimeout(() => {
      tip.innerHTML = '<span class="UnitMetaDetails-toggletipBubble">' + message + '</span>';
    }, 100);
  });
  document.addEventListener('click', e => {
    if (toggle !== e.target) {
      if (!tip) {
        return;
      }
      tip.innerHTML = '';
    }
  });
  toggle.addEventListener('keydown', e => {
    if (!tip) {
      return;
    }
    if (e.key === 'Escape') {
      tip.innerHTML = '';
    }
  });
});
//# sourceMappingURL=toggle-tip.js.map
