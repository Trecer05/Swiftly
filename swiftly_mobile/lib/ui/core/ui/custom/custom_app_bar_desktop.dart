import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:swiftly_mobile/domain/kanban/models/priority.dart';
import 'package:swiftly_mobile/ui/core/themes/colors.dart';
import 'package:swiftly_mobile/ui/core/themes/theme.dart';
import 'package:swiftly_mobile/ui/core/ui/custom/filter_state.dart';

import '../../../../providers/card_notifier_provider.dart';
import '../../../../providers/current_user_provider.dart';
import '../../../../providers/label_notifier_provider.dart';
import '../../../../providers/user_notifier_provider.dart';

class CustomAppBarDesktop extends StatelessWidget {
  final String title;
  final int? quantity;
  final List<Widget> buttons;
  const CustomAppBarDesktop({
    super.key,
    required this.title,
    this.quantity,
    required this.buttons,
  });

  @override
  Widget build(BuildContext context) {
    return Container(
      padding: const EdgeInsets.only(right: 16, left: 16, bottom: 16, top: 48),
      width: double.infinity,
      child: Row(
        children: [
          Stack(
            clipBehavior: Clip.none,
            children: [
              Text(title, style: AppTextStyles.style6),
              if (quantity != null && quantity! > 0)
                Positioned(
                  right: -30,
                  child: Text('($quantity)', style: AppTextStyles.style14),
                ),
            ],
          ),
          const Spacer(),
          ...buttons
              .expand((button) => [button, const SizedBox(width: 5)])
              .toList()
            ..removeLast(),
        ],
      ),
    );
  }
}

class SegmentedControlWidgetDesktop extends ConsumerStatefulWidget {
  const SegmentedControlWidgetDesktop({super.key});

  @override
  ConsumerState<SegmentedControlWidgetDesktop> createState() =>
      _SegmentedControlWidgetDesktopState();
}

class _SegmentedControlWidgetDesktopState
    extends ConsumerState<SegmentedControlWidgetDesktop> {
  final LayerLink _filtersLink = LayerLink();
  OverlayEntry? _overlayEntry;

  bool _menuHovered = false;
  bool _submenuHovered = false;

  void _onMenuHoverChanged(bool hovered) {
    _menuHovered = hovered;
    if (!hovered) _scheduleHideMenu();
  }

  void _onSubmenuHoverChanged(bool hovered) {
    _submenuHovered = hovered;
    if (!hovered) _scheduleHideMenu();
  }

  void _showFiltersMenu() {
    if (_overlayEntry != null) return;

    final overlay = Overlay.of(context);
    _overlayEntry = OverlayEntry(
      builder: (context) {
        return Positioned(
          width: 200,
          child: CompositedTransformFollower(
            link: _filtersLink,
            showWhenUnlinked: false,
            offset: const Offset(0, 40),
            child: MouseRegion(
              cursor: SystemMouseCursors.click,
              onEnter: (_) => _onMenuHoverChanged(true),
              onExit: (_) {
                _onMenuHoverChanged(false);
              },
              child: Container(
                padding: const EdgeInsets.all(4),
                decoration: BoxDecoration(
                  gradient: AppColors.gradient_4,
                  borderRadius: BorderRadius.circular(8),
                ),
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    _FilterMenuItem(
                      title: '---',
                      onSubmenuHoverChanged: _onSubmenuHoverChanged,
                    ),
                    _FilterMenuItem(
                      title: 'Категория',
                      onSubmenuHoverChanged: _onSubmenuHoverChanged,
                    ),
                    _FilterMenuItem(
                      title: 'Приоритет',
                      onSubmenuHoverChanged: _onSubmenuHoverChanged,
                    ),
                    _FilterMenuItem(
                      title: 'Дата',
                      onSubmenuHoverChanged: _onSubmenuHoverChanged,
                    ),
                    _FilterMenuItem(
                      title: 'Исполнитель',
                      onSubmenuHoverChanged: _onSubmenuHoverChanged,
                    ),
                  ],
                ),
              ),
            ),
          ),
        );
      },
    );

    overlay.insert(_overlayEntry!);
  }

  void _hideMenu() {
    _overlayEntry?.remove();
    _overlayEntry = null;
  }

  void _scheduleHideMenu() {
    Future.delayed(const Duration(milliseconds: 180), () {
      if (!_menuHovered && !_submenuHovered) {
        _hideMenu();
      }
    });
  }

  @override
  void dispose() {
    _hideMenu();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    const items = SegmentItem.values;

    return Container(
      padding: const EdgeInsets.all(4),
      height: 40,
      decoration: BoxDecoration(
        color: AppColors.white15,
        borderRadius: BorderRadius.circular(12),
      ),
      child: Row(
        children: [
          for (int i = 0; i < items.length; i++) ...[
            if (items[i].id == 'Filters')
              MouseRegion(
                onEnter: (_) {
                  _onMenuHoverChanged(true);
                  _showFiltersMenu();
                },
                onExit: (_) {
                  _onMenuHoverChanged(false);
                },
                child: CompositedTransformTarget(
                  link: _filtersLink,
                  child: SegmentWidgetDesktop(title: items[i].title),
                ),
              )
            else if (items[i].id == 'My tasks')
              GestureDetector(
                onTap: () {
                  final me = ref.watch(currentUserProvider);
                  if (me != null) {
                    ref.read(filterNotifierProvider.notifier).setUser('${me.name} ${me.lastName ?? ''}');
                  }
                },
                child: SegmentWidgetDesktop(title: items[i].title),
              )
            else
              SegmentWidgetDesktop(title: items[i].title),
            if (i < items.length - 1)
              Container(
                margin: const EdgeInsets.symmetric(horizontal: 8),
                width: 1,
                height: 16,
                color: AppColors.white128,
              ),
          ],
        ],
      ),
    );
  }
}

class _FilterMenuItem extends ConsumerStatefulWidget {
  final String title;
  final void Function(bool) onSubmenuHoverChanged;

  const _FilterMenuItem({
    required this.title,
    required this.onSubmenuHoverChanged,
  });

  @override
  ConsumerState<_FilterMenuItem> createState() => _FilterMenuItemState();
}

class _FilterMenuItemState extends ConsumerState<_FilterMenuItem> {
  bool isHovered = false;

  final LayerLink _submenuLink = LayerLink();
  OverlayEntry? _submenuEntry;
  bool _submenuLocalHovered = false;

  void _showSubmenu() {
    if (_submenuEntry != null) return;

    final overlay = Overlay.of(context)!;

    _submenuEntry = OverlayEntry(
      builder: (context) {
        return Positioned(
          width: 180,
          child: CompositedTransformFollower(
            link: _submenuLink,
            showWhenUnlinked: false,
            offset: const Offset(200, 0),
            child: MouseRegion(
              cursor: SystemMouseCursors.click,
              onEnter: (_) {
                _submenuLocalHovered = true;
                widget.onSubmenuHoverChanged(true);
              },
              onExit: (_) {
                _submenuLocalHovered = false;
                widget.onSubmenuHoverChanged(false);
                Future.delayed(const Duration(milliseconds: 120), () {
                  if (!_submenuLocalHovered) _hideSubmenu();
                });
              },
              child: Container(
                padding: const EdgeInsets.all(4),
                constraints: const BoxConstraints(maxHeight: 300),
                decoration: BoxDecoration(
                  gradient: AppColors.gradient_4,
                  borderRadius: BorderRadius.circular(6),
                ),
                child: SingleChildScrollView(
                  child: Column(
                    crossAxisAlignment: CrossAxisAlignment.start,
                    children: _buildSubmenuOptions(),
                  ),
                ),
              ),
            ),
          ),
        );
      },
    );

    overlay.insert(_submenuEntry!);
  }

  void _hideSubmenu() {
    _submenuEntry?.remove();
    _submenuEntry = null;
  }

  List<Widget> _buildSubmenuOptions() {
    final labels =
        ref
            .watch(labelNotifierProvider)
            .labels
            .map((label) => label.title)
            .toList();
    const priorities = Priority.values;
    final users =
        ref
            .watch(userNotifierProvider)
            .users
            .map((user) => '${user.name} ${user.lastName ?? ''}')
            .toList();
    final List<MapEntry<String, VoidCallback>> options;
    switch (widget.title) {
      case 'Категория':
        final uniqueLabels = {for (var l in labels) l}.toList();
        options = [
          MapEntry('---', () {
            ref.read(filterNotifierProvider.notifier).clearLabel();
          }),
          ...uniqueLabels.map(
            (label) => MapEntry(label, () {
              ref.read(filterNotifierProvider.notifier).setLabel(label);
            }),
          ),
        ];
        break;
      case 'Приоритет':
        options = [
          MapEntry('---', () {
            ref.read(filterNotifierProvider.notifier).clearPriority();
          }),
          ...priorities.map((priority) => MapEntry(priority.title, () {
            ref.read(filterNotifierProvider.notifier).setPriority(priority.title);
          })),
        ];
        break;
      // case 'Дата':
      //   options = ['---', 'Сегодня', 'На этой неделе', 'Выбрать дату'];
      //   break;
      case 'Исполнитель':
        options = [
          MapEntry('---', () {
            ref.read(filterNotifierProvider.notifier).clearUser();
          }),
          ...users.map((user) => MapEntry(user, () {
            ref.read(filterNotifierProvider.notifier).setUser(user);
          })),
        ];
        break;
      default:
        options = [
          MapEntry('Опция 1', () => debugPrint('Опция 1 выбрана')),
          MapEntry('Опция 2', () => debugPrint('Опция 2 выбрана')),
        ];
    }

    return options.map((entry) {
      return SubmenuItemWidget(
        title: entry.key,
        onTap: () {
          entry.value();
          _hideSubmenu();
        },
      );
    }).toList();
  }

  @override
  void dispose() {
    _hideSubmenu();
    super.dispose();
  }

  void _onEnter() {
    setState(() => isHovered = true);
    _showSubmenu();
    widget.onSubmenuHoverChanged(true);
  }

  void _onExit() {
    setState(() => isHovered = false);
    Future.delayed(const Duration(milliseconds: 120), () {
      if (!_submenuLocalHovered) {
        _hideSubmenu();
        widget.onSubmenuHoverChanged(false);
      }
    });
  }

  @override
  Widget build(BuildContext context) {
    return MouseRegion(
      onEnter: (_) => _onEnter(),
      onExit: (_) => _onExit(),
      child: CompositedTransformTarget(
        link: _submenuLink,
        child: Container(
          width: double.infinity,
          padding: const EdgeInsets.symmetric(horizontal: 12, vertical: 8),
          decoration: BoxDecoration(
            color: isHovered ? AppColors.white31 : Colors.transparent,
            borderRadius: BorderRadius.circular(6),
          ),
          child: Row(
            children: [
              Expanded(child: Text(widget.title, style: AppTextStyles.style17)),
              const Icon(Icons.arrow_right, size: 16, color: AppColors.white),
            ],
          ),
        ),
      ),
    );
  }
}

class SubmenuItemWidget extends StatefulWidget {
  final String title;
  final VoidCallback? onTap;
  const SubmenuItemWidget({super.key, required this.title, this.onTap});

  @override
  State<SubmenuItemWidget> createState() => _SubmenuItemWidgetState();
}

class _SubmenuItemWidgetState extends State<SubmenuItemWidget> {
  bool isHovered = false;

  @override
  Widget build(BuildContext context) {
    return MouseRegion(
      onEnter:
          (_) => setState(
            () => isHovered = widget.title != 'Пусто' ? true : false,
          ),
      onExit: (_) => setState(() => isHovered = false),
      child: GestureDetector(
        onTap: widget.title == 'Пусто' ? null : widget.onTap,
        child: Container(
          width: double.infinity,
          padding: const EdgeInsets.symmetric(horizontal: 12, vertical: 8),
          decoration: BoxDecoration(
            color: isHovered ? AppColors.white31 : Colors.transparent,
            borderRadius: BorderRadius.circular(4),
          ),
          child: Text(widget.title, style: AppTextStyles.style17),
        ),
      ),
    );
  }
}

class SegmentWidgetDesktop extends StatefulWidget {
  final String title;
  const SegmentWidgetDesktop({super.key, required this.title});

  @override
  State<SegmentWidgetDesktop> createState() => _SegmentWidgetDesktopState();
}

class _SegmentWidgetDesktopState extends State<SegmentWidgetDesktop> {
  bool isHovered = false;

  @override
  Widget build(BuildContext context) {
    return MouseRegion(
      cursor: SystemMouseCursors.click,
      onEnter: (_) => setState(() => isHovered = true),
      onExit: (_) => setState(() => isHovered = false),
      child: Container(
        padding: const EdgeInsets.symmetric(horizontal: 12, vertical: 8),
        decoration: BoxDecoration(
          color: isHovered ? AppColors.white15 : null,
          borderRadius: BorderRadius.circular(8),
        ),
        child: Text(
          widget.title,
          style: AppTextStyles.style17.copyWith(height: 1.0),
          overflow: TextOverflow.ellipsis,
        ),
      ),
    );
  }
}

enum SegmentItem {
  filters('Filters', 'Фильтры'),
  sort('Sort', 'Сортировать'),
  myTasks('My tasks', 'Мои задачи');

  final String id;
  final String title;

  const SegmentItem(this.id, this.title);
}
