import 'package:flutter/material.dart';

import '../themes/colors.dart';

class NavigationState {
  final bool showLabels;
  final int? selectedIndex;

  NavigationState({required this.showLabels, this.selectedIndex});
}

class CustomNavigationRail extends StatefulWidget {
  final ValueNotifier<NavigationState> stateNotifier =
      ValueNotifier<NavigationState>(
        NavigationState(showLabels: true, selectedIndex: 0),
      );
  late final List<ItemWidget> itemWidget;
  CustomNavigationRail({super.key}) {
    itemWidget = [
      ItemWidget(
        iconData: Icons.sort,
        isController: true,
        stateNotifier: stateNotifier,
      ),
      ItemWidget(
        iconData: Icons.home,
        label: 'Главная',
        index: 0,
        isController: false,
        stateNotifier: stateNotifier,
      ),
      ItemWidget(
        iconData: Icons.chat,
        label: 'Чат',
        index: 1,
        isController: false,
        stateNotifier: stateNotifier,
      ),
      ItemWidget(
        iconData: Icons.code,
        label: 'Код',
        index: 2,
        isController: false,
        stateNotifier: stateNotifier,
      ),
      ItemWidget(
        iconData: Icons.cloud,
        label: 'Облако',
        index: 3,
        isController: false,
        stateNotifier: stateNotifier,
      ),
      ItemWidget(
        iconData: Icons.pan_tool,
        label: 'Фигма',
        index: 4,
        isController: false,
        stateNotifier: stateNotifier,
      ),
      ItemWidget(
        iconData: Icons.task,
        label: 'Задачи',
        index: 5,
        isController: false,
        stateNotifier: stateNotifier,
      ),
      ItemWidget(
        iconData: Icons.settings,
        label: 'Settings',
        index: 6,
        isController: false,
        stateNotifier: stateNotifier,
      ),
    ];
  }

  @override
  State<CustomNavigationRail> createState() => _CustomNavigationRailState();
}

class _CustomNavigationRailState extends State<CustomNavigationRail> {
  final double expandedWidth = 200;
  final double collapsedWidth = 60;

  @override
  void dispose() {
    widget.stateNotifier.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      body: Row(
        mainAxisAlignment: MainAxisAlignment.start,
        children: [
          ValueListenableBuilder<NavigationState>(
            valueListenable: widget.stateNotifier,
            builder: (_, state, __) {
              return Container(
                padding: EdgeInsets.symmetric(vertical: 10, horizontal: 10),
                width: state.showLabels ? expandedWidth : null,
                decoration: BoxDecoration(
                  color: const Color.fromARGB(255, 3, 13, 55),
                ),
                child: IntrinsicWidth(
                  child: Column(
                    mainAxisAlignment: MainAxisAlignment.start,
                    crossAxisAlignment: CrossAxisAlignment.start,
                    children: [
                      for (int i = 0; i < widget.itemWidget.length - 1; i++)
                        widget.itemWidget[i],
                      const Spacer(),
                      widget.itemWidget.last,
                    ],
                  ),
                ),
              );
            },
          ),
        ],
      ),
    );
  }
}

class ItemWidget extends StatelessWidget {
  final IconData iconData;
  final String? label;
  final int? index;
  final ValueNotifier<NavigationState> stateNotifier;
  final bool isController;

  const ItemWidget({
    super.key,
    required this.iconData,
    this.label,
    this.index,
    required this.stateNotifier,
    required this.isController,
  });

  @override
  Widget build(BuildContext context) {
    return ValueListenableBuilder<NavigationState>(
      valueListenable: stateNotifier,
      builder: (_, state, __) {
        final isSelected =
            !isController && index != null && state.selectedIndex == index;

        return GestureDetector(
          onTap: () {
            if (isController) {
              stateNotifier.value = NavigationState(
                showLabels: !state.showLabels,
                selectedIndex: state.selectedIndex,
              );
            } else if (index != null) {
              stateNotifier.value = NavigationState(
                showLabels: state.showLabels,
                selectedIndex: index,
              );
            }
          },
          child: Container(
            padding: const EdgeInsets.all(8),
            width: double.infinity,
            decoration: BoxDecoration(
              gradient: isSelected ? AppColors.gradient_4 : null,
              borderRadius: BorderRadius.circular(12),
            ),
            child: Row(
              mainAxisSize: MainAxisSize.min,
              children: [
                Icon(
                  iconData,
                  color: isSelected ? AppColors.white : AppColors.white128,
                ),
                if (state.showLabels && label != null) ...[
                  const SizedBox(width: 8),
                  Text(
                    label!,
                    style: TextStyle(
                      color: isSelected ? AppColors.white : AppColors.white128,
                    ),
                  ),
                ],
              ],
            ),
          ),
        );
      },
    );
  }
}